export class Exercise {
    Id: number;
    Classid: any;
    Title: string;
    Description: string;
    Link: string;
    Status: any;
    Create_date: any;
}

export class Submission {
    Id: number;
    Studentid: number;
    Exerciseid: number;
    Description: string;
    Link: string;
    Status: any;
    Create_date: any;
}
