export class User {
    Id: number;
    Username: string;
    Password: string;
    Firstname: string;
    Lastname: string;
    authdata?: string;
}

export class UserBackend {
    Id: number;
    Username: string;
    Password: string;
    Firstname: string;
    Lastname: string;
    Email: string;
    Tel: string;
    Avatar: string;
    Role: any;
    Status: any;
    Create_date: any;
}

export class ServerResponse {
    result: any;
}
