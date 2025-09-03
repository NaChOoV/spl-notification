import axios, { type AxiosInstance } from 'axios';
import EnvConfig from '../config/enviroment';
import type { AbmUser, User } from '../types/user.type';

class SourceService {
    private readonly httpService: AxiosInstance;
    constructor() {
        this.httpService = axios.create({
            baseURL: EnvConfig.sourceServiceBaseUrl,
            timeout: 10000,
            headers: {
                'X-Auth-String': EnvConfig.sourceServiceAuthString,
            },
        });
    }

    public async getAbmUserByRun(run: string): Promise<AbmUser | undefined> {
        const response = await this.httpService.get<AbmUser | undefined>(`/user/abm/${run}`);
        return response.data;
    }

    public async getUserByExternalId(externalId: number): Promise<User | undefined> {
        const response = await this.httpService.get<User | undefined>(`/user/${externalId}`);
        return response.data;
    }
}

export default SourceService;
