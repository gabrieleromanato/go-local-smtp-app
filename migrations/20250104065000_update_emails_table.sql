UPDATE emails
SET body_html = ''
WHERE body_html IS NULL;