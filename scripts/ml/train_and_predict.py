import os
import time
import psycopg2
import pandas as pd
from sklearn.ensemble import RandomForestClassifier
from sklearn.preprocessing import LabelEncoder
from dotenv import load_dotenv

def wait_for_db(conn_kwargs, max_retries=10, delay=5):
    print(f"Waiting for database at {conn_kwargs['host']}:{conn_kwargs['port']}...")
    for i in range(max_retries):
        try:
            conn = psycopg2.connect(**conn_kwargs)
            conn.close()
            print("Database is ready!")
            return True
        except Exception as e:
            print(f"Database not ready (attempt {i+1}/{max_retries}): {e}")
            time.sleep(delay)
    return False

def main():
    load_dotenv('../../.env', override=True)
    db_host = os.environ.get('DB_HOST_OVERRIDE', os.getenv('DB_HOST', 'localhost'))
    if db_host == 'postgres' and not os.path.exists('/.dockerenv'):
        db_host = 'localhost'

    conn_kwargs = {
        'dbname': os.getenv('DB_NAME', 'upcycleconnect'),
        'user': os.getenv('DB_USER', 'admin'),
        'password': os.getenv('DB_PASSWORD', 'root'),
        'host': db_host,
        'port': os.getenv('DB_PORT', '5432')
    }

    if not wait_for_db(conn_kwargs):
        print("Could not connect to database. Exiting.")
        return

    print("Fetching data from database...")
    conn = psycopg2.connect(**conn_kwargs)
    query = """
    SELECT 
        u.id as user_id,
        u.role,
        EXTRACT(DAY FROM (NOW() - u.created_at)) as days_since_registration,
        COALESCE((SELECT SUM(points) FROM score_history WHERE user_id = u.id), 0) as total_points,
        (SELECT COUNT(*) FROM event_participation WHERE user_id = u.id) as count_events,
        (SELECT COUNT(*) FROM course_order WHERE buyer_id = u.id) as count_courses,
        (SELECT COUNT(*) FROM listing_order WHERE user_id = u.id) as count_listings
    FROM users u;
    """
    import warnings
    with warnings.catch_warnings():
        warnings.simplefilter('ignore', UserWarning)
        df = pd.read_sql_query(query, conn)

    if len(df) < 10:
        print("Not enough data to train.")
        conn.close()
        return

    print("Feature Engineering...")
    le = LabelEncoder()
    df['role_encoded'] = le.fit_transform(df['role'])
    
    X = df[['role_encoded', 'days_since_registration', 'total_points']]
    
    def determine_target(row):
        counts = {'event': row['count_events'], 'course': row['count_courses'], 'listing': row['count_listings']}
        max_service = max(counts, key=counts.get)
        return max_service if counts[max_service] > 0 else 'listing'

    df['target'] = df.apply(determine_target, axis=1)
    y = df['target']

    print("Training Random Forest model...")
    model = RandomForestClassifier(n_estimators=100, random_state=42)
    model.fit(X, y)

    print("Predicting for all users...")
    df['predicted_service'] = model.predict(X)
    df['probability'] = model.predict_proba(X).max(axis=1)

    print("Saving predictions to database...")
    cur = conn.cursor()
    for index, row in df.iterrows():
        cur.execute("""
            INSERT INTO user_predictions (user_id, predicted_service_type, probability, calculated_at)
            VALUES (%s, %s, %s, NOW())
            ON CONFLICT (user_id) DO UPDATE 
            SET predicted_service_type = EXCLUDED.predicted_service_type,
                probability = EXCLUDED.probability,
                calculated_at = NOW();
        """, (int(row['user_id']), row['predicted_service'], float(row['probability'])))
    conn.commit()
    cur.close()
    conn.close()
    print("Done!")

if __name__ == "__main__":
    ONE_HOUR = 3600
    while True:
        main()
        time.sleep(ONE_HOUR)

