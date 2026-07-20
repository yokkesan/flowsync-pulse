export type LoginRequest = {
    email: string;
    password: string;
};

export type AuthUser = {
    id: number;
    display_name: string;
    email: string;
    role: string;
};

export type AuthCompany = {
    id: number;
    name: string;
};

export type LoginResponse = {
    access_token: string;
    token_type: 'Bearer';
    expires_in: number;
    user: AuthUser;
    company: AuthCompany;
};

export type CurrentUserResponse = {
    user: AuthUser;
    company: AuthCompany;
};