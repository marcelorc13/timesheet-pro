-- +goose Up
-- +goose StatementBegin
CREATE TABLE timesheet_statuses (
  id SMALLINT PRIMARY KEY,
  name TEXT NOT NULL
);
INSERT INTO timesheet_statuses VALUES (1, 'open'), (2, 'closed'), (3, 'absent'), (4, 'approved'), (5, 'reproved') ;

CREATE TABLE daily_timesheets (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id),
  organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  date DATE NOT NULL, 
  status_id SMALLINT NOT NULL REFERENCES timesheet_statuses(id),
  total_minutes BIGINT DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(user_id, organization_id, date)
);

CREATE TABLE entry_types (
  id SMALLINT PRIMARY KEY,
  name TEXT NOT NULL
);
INSERT INTO entry_types VALUES (1, 'in'), (2, 'out');

CREATE TABLE timesheet_entries (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  timesheet_id UUID NOT NULL REFERENCES daily_timesheets(id) ON DELETE CASCADE,
  organization_id UUID NOT NULL REFERENCES organizations(id),
  type_id SMALLINT NOT NULL REFERENCES entry_types(id),
  timestamp TIMESTAMPTZ DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE timesheet_entries;
DROP TABLE entry_types;
DROP TABLE daily_timesheets;
DROP TABLE timesheet_statuses;
-- +goose StatementEnd
