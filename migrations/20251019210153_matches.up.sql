CREATE TABLE matches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    freelancer_id UUID NOT NULL REFERENCES freelancers(id) ON DELETE CASCADE,
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT now(),
    UNIQUE (freelancer_id, project_id)
);