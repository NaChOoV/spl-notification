import type { AxiosInstance } from 'axios';
import axios from 'axios';
import EnvConfig from '../config/enviroment';
import type { NotifyTrack } from '../types/notify-track';
import { getLocation } from '../utils/location';
import type { Track } from '../db/schema';

type NotifyResult = {
    fullFilled: number;
    rejected: number;
};

class WhatsappService {
    private readonly httpService: AxiosInstance;
    constructor() {
        this.httpService = axios.create({
            baseURL: EnvConfig.whatsappBaseUrl,
            timeout: 10000,
            auth: {
                username: EnvConfig.whatsappUsername,
                password: EnvConfig.whatsappPassword,
            },
        });
    }

    public async notifyEntry(tracks: NotifyTrack[]): Promise<NotifyResult> {
        const promises = tracks.map<Promise<void>>((track) => {
            const body = {
                chatId: String(track.chatId),
                fullName: track.fullName,
                location: getLocation(track.location),
            };

            return this.httpService.post('webhook/whatsapp/notify-entry', body);
        });

        const response = await Promise.allSettled(promises);
        const fulfilled = response.filter((res) => res.status === 'fulfilled');
        const rejected = response.filter((res) => res.status === 'rejected');

        return {
            fullFilled: fulfilled.length,
            rejected: rejected.length,
        };
    }

    public async listTracks(chatId: string, tracks: Track[]): Promise<void> {
        let message = 'No tienes seguimientos.';
        if (tracks.length > 0) {
            message = `📋 Listado:\n${tracks
                .map((track) => `- ${track.run} ${track.alias || track.fullName}`)
                .join('\n')}`;
        }
        return this.sendMessage(chatId, message);
    }

    public sendMessage(chatId: string, message: string): Promise<void> {
        return this.httpService.post('webhook/whatsapp', { chatId, message });
    }
}

export default WhatsappService;
