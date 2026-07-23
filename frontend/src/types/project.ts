export type ProjectMember = {
    user_id: number;
    display_name: string;
    role: string;
    status: string;
};

export type Project = {
    project_id: number;
    company_id: number;
    name: string;
    slug: string;
    description: string;
    status: string;
    start_date: string | null;
    end_date: string | null;
    task_count: number;
    members: ProjectMember[];
    created_at: string;
    updated_at: string;
};

export type ProjectListResponse = {
    projects: Project[];
};

export type ProjectStatusFilter =
    | 'all'
    | 'active'
    | 'completed'
    | 'paused';