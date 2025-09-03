import type { AxiosInstance } from 'axios';
import { TrackType, type Track, type TrackWithLocation } from '../db/schema';
import type TrackRepository from '../repository/track.repository';
import type { Access } from '../types/access';
import type { NotifyTrack } from '../types/notify-track';
import { extractUnique } from '../utils/utility';
import type WhatsappService from './whatsapp.service';
import axios from 'axios';
import EnvConfig from '../config/enviroment';

class AccessService {
    private readonly httpService: AxiosInstance;
    constructor(
        private readonly trackRepository: TrackRepository,
        private readonly whatsappService: WhatsappService
    ) {
        this.httpService = axios.create({
            baseURL: `${EnvConfig.accessServiceBaseUrl}/api`,
            timeout: 10000,
            headers: {
                'X-Auth-Token': EnvConfig.accessServiceAuthString,
            },
        });
    }

    public async checkAccess(accessArray: Access[]): Promise<void> {
        if (accessArray.length === 0) return;

        const matchTracks: Track[] = [];
        const allTracks = await this.trackRepository.getAll(TrackType.TRACK);

        allTracks.forEach((track) => {
            const matchAccess = accessArray.find(
                (access) => access.externalId === track.userId && access.entryAt !== track.lastEntry
            );
            if (!matchAccess) return;

            matchTracks.push(track);
        });
        if (matchTracks.length === 0) return;

        const matchUserIds = extractUnique(matchTracks, (track) => track.userId);
        const accessToUpdate = accessArray.filter((access) =>
            matchUserIds.includes(access.externalId)
        );

        await this.trackRepository.updateTrack(accessToUpdate, TrackType.TRACK);

        const tracksToNotify: NotifyTrack[] = [];
        matchTracks.forEach((track) => {
            const access = accessToUpdate.find((access) => access.externalId === track.userId);
            if (!access) return;

            tracksToNotify.push({
                chatId: track.chatId,
                run: track.run,
                fullName: track.fullName,
                alias: track.alias,
                location: access.location,
            });
        });

        console.log('Nuevos seguimientos:', tracksToNotify.length);
        const result = await this.whatsappService.notifyEntry(tracksToNotify);
        console.log(`FullFilled: ${result.fullFilled}, Rejected: ${result.rejected}`);
    }

    public async getAccess(): Promise<Access[]> {
        const response = await this.httpService.get<{ data: Access[] }>('/access/complete');
        return response.data.data;
    }
}

export default AccessService;
