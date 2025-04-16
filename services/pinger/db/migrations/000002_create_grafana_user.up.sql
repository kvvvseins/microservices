CREATE USER grafanareader WITH PASSWORD 'grafanareader';
ALTER USER grafanareader WITH PASSWORD 'grafanareader';
GRANT CONNECT ON DATABASE pinger TO grafanareader;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO grafanareader;