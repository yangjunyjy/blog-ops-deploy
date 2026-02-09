```sql
#查询长事务
SELECT 
    pid, 
    usename, 
    application_name, 
    client_addr, 
    backend_start, 
    query_start, 
    now() - query_start AS duration, 
    state, 
    query
FROM 
    pg_stat_activity
WHERE 
    state != 'idle' 
    AND now() - query_start > interval '5 minutes'
ORDER BY 
    duration DESC;
    
#杀会话
SELECT pg_terminate_backend(pid);
```

