import TrackRepository from './repository/track.repository';
import AccessService from './service/access.service';
import SourceService from './service/source.service';
import WhatsappService from './service/whatsapp.service';

const whatsappService = new WhatsappService();
const trackRepository = new TrackRepository();
const sourceService = new SourceService();

const accessService = new AccessService(trackRepository, whatsappService);

export { whatsappService, trackRepository, accessService, sourceService };
