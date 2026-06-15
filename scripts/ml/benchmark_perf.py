import os
import time
import psycopg2
from dotenv import load_dotenv

def benchmark_query(cur, name, query, params=None):
    start = time.time()
    cur.execute(f"EXPLAIN ANALYZE {query}", params)
    analysis = cur.fetchall()
    end = time.time()
    
    # Extract execution time from EXPLAIN ANALYZE
    exec_time = (end - start) * 1000
    print(f"Query: {name}")
    print(f"Execution Time: {exec_time:.2f} ms")
    # print("\n".join([line[0] for line in analysis]))
    print("-" * 30)
    return exec_time

def main():
    load_dotenv('.env')
    conn = psycopg2.connect(
        dbname=os.getenv('DB_NAME', 'upcycleconnect'),
        user=os.getenv('DB_USER', 'admin'),
        password=os.getenv('DB_PASSWORD', 'root'),
        host=os.getenv('DB_HOST', 'localhost'),
        port=os.getenv('DB_PORT', '5432')
    )
    cur = conn.cursor()

    print("--- BENCHMARK MODULE 2 (15,000 users simulation) ---")

    # 1. Query: Get Score History for a specific user
    # Problem: No index on score_history(user_id)
    q1 = "SELECT * FROM score_history WHERE user_id = %s"
    benchmark_query(cur, "Score History (Without Index)", q1, (1,))

    # 2. Query: Global Reporting Stats
    q2 = """
        SELECT 'event' as type, COUNT(*), SUM(price) FROM event e 
        JOIN event_participation ep ON e.id = ep.event_id GROUP BY type
        UNION ALL
        SELECT 'course' as type, COUNT(*), SUM(price) FROM course_order GROUP BY type
    """
    benchmark_query(cur, "Global Stats (Complex Join)", q2)

    print("\nApplying Optimizations...\n")
    cur.execute("CREATE INDEX IF NOT EXISTS idx_score_history_user_id ON score_history(user_id);")
    conn.commit()

    print("--- POST-OPTIMIZATION ---")
    benchmark_query(cur, "Score History (With Index)", q1, (1,))

    cur.close()
    conn.close()

if __name__ == "__main__":
    main()
