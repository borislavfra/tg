import {SCHEMAS} from '../_schemas';
import {RequestParamsType} from './types';

const ENDPOINT = '/adder/getUUID'

export const makeRequestConfig = ({additionalFetchParams,bodyParams}: RequestParamsType)=>({endpoint: ENDPOINT,responseSchema: SCHEMAS.getUUID,body: {params: bodyParams},...additionalFetchParams});