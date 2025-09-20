import type { AxiosInstance } from 'axios';
import { TrackType, type Track, type TrackWithLocation } from '../db/schema';
import type TrackRepository from '../repository/track.repository';
import type { Access } from '../types/access';
import type { NotifyTrack } from '../types/notify-track';
import { extractUnique } from '../utils/utility';
import type WhatsappService from './whatsapp.service';
import axios from 'axios';
import EnvConfig from '../config/enviroment';
import { Agent } from 'https';

class AccessService {
    private readonly httpService: AxiosInstance;
    constructor(
        private readonly trackRepository: TrackRepository,
        private readonly whatsappService: WhatsappService
    ) {
        this.httpService = axios.create({
            baseURL: `${EnvConfig.accessServiceBaseUrl}/api`,
            timeout: 10000,
            httpsAgent: new Agent({
                keepAlive: true,
                keepAliveMsecs: 600000,
                maxSockets: 3,
                maxFreeSockets: 2,
                maxTotalSockets: 5,
            }),
            headers: {
                Connection: 'keep-alive',
                'Keep-Alive': 'timeout=600, max=6000',
                'X-Auth-Token': EnvConfig.accessServiceAuthString,
            },
        });
    }

    public async checkAccess(accessArray: Access[]): Promise<void> {
        let allTracks = await this.trackRepository.getAll(TrackType.TRACK);

        const matchEntryAtTracks = allTracks.filter((track) => {
            return accessArray.some(
                (access) => access.externalId === track.userId && access.entryAt !== track.lastEntry
            );
        }) as Track[];

        const matchExitAtTracks = allTracks.filter((track) => {
            return accessArray.some(
                (access) =>
                    access.externalId === track.userId &&
                    access.exitAt &&
                    access.exitAt !== track.lastExit
            );
        }) as Track[];

        const matchEntryUserIds = extractUnique(matchEntryAtTracks, (track) => track.userId);
        const accessToUpdateEntry = accessArray.filter((access) =>
            matchEntryUserIds.includes(access.externalId)
        );

        const matchExitUserIds = extractUnique(matchExitAtTracks, (track) => track.userId);
        const accessToUpdateExit = accessArray.filter((access) =>
            matchExitUserIds.includes(access.externalId)
        );

        if (accessToUpdateEntry.length !== 0) {
            await this.trackRepository.updateEntryAt(accessToUpdateEntry, TrackType.TRACK);
        }
        if (accessToUpdateExit.length !== 0) {
            await this.trackRepository.updateExitAt(accessToUpdateExit, TrackType.TRACK);
        }

        const tracksToNotifyEntry: NotifyTrack[] = [];
        matchEntryAtTracks.forEach((track) => {
            const access = accessToUpdateEntry.find((access) => access.externalId === track.userId);
            if (!access) return;

            tracksToNotifyEntry.push({
                chatId: track.chatId,
                run: track.run,
                fullName: track.fullName,
                alias: track.alias,
                location: access.location,
            });
        });

        const tracksToNotifyExit: NotifyTrack[] = [];
        matchExitAtTracks.forEach((track) => {
            const access = accessToUpdateExit.find((access) => access.externalId === track.userId);
            if (!access) return;

            tracksToNotifyExit.push({
                chatId: track.chatId,
                run: track.run,
                fullName: track.fullName,
                alias: track.alias,
                location: access.location,
            });
        });

        if (tracksToNotifyEntry.length > 0) {
            const result = await this.whatsappService.notifyEntry(tracksToNotifyEntry);
            console.log('Entradas registradas:', tracksToNotifyEntry.length);
            console.log(
                `FullFilled Entry: ${result.fullFilled}, Rejected Entry: ${result.rejected}`
            );
        }

        if (tracksToNotifyExit.length > 0) {
            const result = await this.whatsappService.notifyExit(tracksToNotifyExit);
            console.log('Salidas registradas:', tracksToNotifyExit.length);
            console.log(`FullFilled Exit: ${result.fullFilled}, Rejected Exit: ${result.rejected}`);
        }

        matchEntryAtTracks.length = 0;
        matchExitAtTracks.length = 0;
        allTracks.length = 0;

        allTracks = null as any;
    }

    public async getAccess(): Promise<Access[]> {
        const response = await this.httpService.get<{ data: Access[] }>('/access/complete');
        return response.data.data;
    }
}

export default AccessService;
