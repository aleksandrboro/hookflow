CREATE TABLE IF NOT EXISTS jobs (
    id UUID PRIMARY KEY,
    target_url TEXT NOT NULL,
    payload JSONB NOT NULL CHECK (jsonb_typeof(payload) = 'object'),
    status TEXT NOT NULL DEFAULT 'pending' CHECK (
        status IN (
            'pending',
            'queued',
            'running',
            'retrying',
            'succeeded',
            'failed',
            'dead',
            'cancelled'
        )
    ),
    run_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    attempts_count INTEGER NOT NULL DEFAULT 0 CHECK (attempts_count >= 0),
    max_attempts INTEGER NOT NULL DEFAULT 5 CHECK (
        max_attempts BETWEEN 1 AND 10
    ),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);