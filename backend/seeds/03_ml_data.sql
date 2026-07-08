




DO $$
DECLARE
    i INT;
    new_user_id BIGINT;
    service_roll INT;
    interaction_count INT;
    item_count INT;
    comment_count INT;
    msg_count INT;
BEGIN
    FOR i IN 1..600 LOOP
        
        INSERT INTO users (username, first_name, last_name, email, password_hash, role, created_at)
        VALUES (
            'trash_user_' || i, 
            'TrashFirst' || i, 
            'TrashLast' || i, 
            'trash_' || i || '@upcycleconnect.test', 
            '$2b$10$hashhashhashhashhashhashhashhashhashhash', 
            CASE 
                WHEN (i % 7 = 0) THEN 'pro'::USER_ROLE 
                WHEN (i % 13 = 0) THEN 'interne'::USER_ROLE
                ELSE 'client'::USER_ROLE 
            END,
            NOW() - (random() * INTERVAL '500 days')
        ) RETURNING id INTO new_user_id;

        
        INSERT INTO score_history (user_id, points, description, created_at)
        SELECT new_user_id, floor(random() * 50), 'Points automatiques ' || s, NOW() - (random() * INTERVAL '365 days')
        FROM generate_series(1, (1 + floor(random() * 5))::INT) s;

        
        
        
        
        
        
        interaction_count := 2 + floor(random() * 10);
        
        FOR j IN 1..interaction_count LOOP
            IF (i % 3 = 0) THEN
                
                service_roll := CASE WHEN random() < 0.7 THEN 0 ELSE floor(random() * 3) END;
            ELSIF (i % 3 = 1) THEN
                
                service_roll := CASE WHEN random() < 0.7 THEN 1 ELSE floor(random() * 3) END;
            ELSE
                
                service_roll := CASE WHEN random() < 0.7 THEN 2 ELSE floor(random() * 3) END;
            END IF;

            IF service_roll = 0 THEN
                INSERT INTO event_participation (event_id, user_id, stripe_payment_intent_id)
                VALUES (1, new_user_id, 'pi_ml_ev_' || i || '_' || j) ON CONFLICT DO NOTHING;
            ELSIF service_roll = 1 THEN
                INSERT INTO course_order (course_id, buyer_id, price, booked_at, stripe_payment_intent_id)
                VALUES (1, new_user_id, 99.99, NOW() - (random() * INTERVAL '60 days'), 'pi_ml_co_' || i || '_' || j);
            ELSE
                INSERT INTO listing_order (listing_id, user_id, price, created_at, status, stripe_payment_intent_id)
                VALUES ((1 + floor(random() * 2)), new_user_id, 150.00, NOW() - (random() * INTERVAL '60 days'), 'completed', 'pi_ml_li_' || i || '_' || j);
            END IF;
        END LOOP;

        
        item_count := floor(random() * 8);
        IF item_count > 0 THEN
            INSERT INTO item (owner_id, site_id, material_type, physical_state, status, created_at)
            SELECT new_user_id, 1, 
                (ARRAY['Bois', 'Métal', 'Plastique', 'Tissu'])[(floor(random() * 4) + 1)::INT],
                (ARRAY['neuf'::ITEM_STATE, 'bon etat'::ITEM_STATE, 'abime'::ITEM_STATE, 'casse'::ITEM_STATE])[(floor(random() * 4) + 1)::INT],
                'collected',
                NOW() - (random() * INTERVAL '100 days')
            FROM generate_series(1, item_count);
        END IF;

        
        comment_count := floor(random() * 4);
        IF comment_count > 0 THEN
            INSERT INTO comments (news_id, created_by, content, created_at)
            SELECT 1, new_user_id, 'Commentaire ML ' || i || '-' || s, NOW() - (random() * INTERVAL '50 days')
            FROM generate_series(1, comment_count) s;
        END IF;

        
        msg_count := floor(random() * 6);
        IF msg_count > 0 THEN
            INSERT INTO chat_message (conversation_id, sender_id, content, created_at)
            SELECT 1, new_user_id, 'Chat ML ' || i || '-' || s, NOW() - (random() * INTERVAL '50 days')
            FROM generate_series(1, msg_count) s;
        END IF;
    END LOOP;
END $$;
