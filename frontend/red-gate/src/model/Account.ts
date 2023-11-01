interface Account {
    account_id: string | null;
    username: string;
    email: string;
    passwordHash: string;
    createdAt: string | null;  // Assuming this is a string in ISO 8601 format
    isSubscribe: boolean | null;
}