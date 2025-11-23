
CREATE TYPE strike_reason_type AS ENUM ('spam', 'fake_profile', 'offensive_content', 'harassment', 'inappropriate_content', 'underage', 'copyright_violation', 'other');
CREATE TYPE strike_status_type AS ENUM ('pending', 'approved', 'rejected', 'resolved');

CREATE TABLE strike (
                         id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                         reporter_id UUID NOT NULL REFERENCES "user" (id) ON DELETE CASCADE,
                         target_user_id UUID NOT NULL REFERENCES "user" (id) ON DELETE CASCADE,
                         type strike_reason_type NOT NULL,
                         reason TEXT,
                         status strike_status_type NOT NULL DEFAULT 'pending',
                         created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                         updated_at TIMESTAMPTZ,
                         moderator_id UUID REFERENCES "user" (id) ON DELETE CASCADE,
                         moderator_note TEXT,

                         UNIQUE (reporter_id, target_user_id),
                         CONSTRAINT strike_reason_check CHECK (LENGTH(reason) BETWEEN 1 AND 250),
                         CONSTRAINT strike_note_check CHECK (LENGTH(moderator_note) BETWEEN 1 AND 250),
                         CONSTRAINT chk_strikes_dates CHECK (created_at <= COALESCE(updated_at, NOW())),
                         CONSTRAINT chk_strikes_self_report CHECK (reporter_id != target_user_id)
);