CREATE ROLE readonly_interviewers;

GRANT CONNECT ON DATABASE postgres TO readonly_interviewers;
GRANT USAGE ON SCHEMA public TO readonly_interviewers;

GRANT SELECT ON members TO readonly_interviewers;
GRANT SELECT ON scores TO readonly_interviewers;

CREATE USER interviewer WITH PASSWORD '';

GRANT readonly_interviewers TO interviewer;