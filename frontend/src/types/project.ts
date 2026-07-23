export type ProjectStatus =
    | 'planned'
    | 'active'
    | 'completed'
    | 'archived';

export type ProjectStatusFilter =
    | 'all'
    | ProjectStatus;

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
    description: string | null;
    status: ProjectStatus;
    start_date: string | null;
    end_date: string | null;
    members: ProjectMember[];
    task_count: number;
    created_at: string;
    updated_at: string;
};

export type ProjectListResponse = {
    projects: Project[];
};

export type ProjectWriteRequest = {
    name: string;
    slug: string;
    description: string | null;
    status: ProjectStatus;
    start_date: string | null;
    end_date: string | null;
    member_ids: number[];
};

export type CreateProjectRequest =
    ProjectWriteRequest;

export type UpdateProjectRequest =
    ProjectWriteRequest;