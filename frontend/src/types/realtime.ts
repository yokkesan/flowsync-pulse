export type WorkContextChangedPayload = {
    user_id: number;
    project_id: number | null;
    project_name: string | null;
    task_id: number | null;
    task_key: string | null;
    task_name: string | null;
    repository_name: string | null;
    branch_name: string | null;
    ticket_key: string | null;
    workspace_name: string | null;
    match_status: string | null;
    session_status: string | null;
    extension_active: boolean;
    started_at: string | null;
    last_heartbeat_at: string | null;
    ended_at: string | null;
    end_reason: string | null;
};

export type RealtimeEvent =
    | {
        type: 'work_context.changed';
        payload: WorkContextChangedPayload;
    };