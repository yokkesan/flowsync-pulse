export type CompanyUser = {
    user_id: number;
    display_name: string;
    email: string;
    role: string;
    status: string;
};

export type CompanyUserListResponse = {
    users: CompanyUser[];
};