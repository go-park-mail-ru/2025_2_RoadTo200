SELECT column_name 
FROM information_schema.columns 
WHERE table_name = 'user' 
ORDER BY ordinal_position;