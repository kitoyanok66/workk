CREATE TABLE freelancer_skills (
    freelancer_id UUID REFERENCES freelancers(id) ON DELETE CASCADE,
    skill_id UUID REFERENCES skills(id) ON DELETE CASCADE,
    PRIMARY KEY (freelancer_id, skill_id)
);