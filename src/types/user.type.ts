export type AbmUser = {
    externalId: number;
    run: string;
    firstName: string;
    lastName: string;
};

export type User = {
    run: string;
    firstName: string;
    lastName: string;
    imageUrl: string;
    access: {
        location: number;
        entryAt: string;
        exitAt: string;
    }[];
};

export enum State {
    ACTIVO = 1,
    INACTIVO = 0,
}

export enum AccessType {
    BIOMETRIA = 1,
    QR = 2,
}
