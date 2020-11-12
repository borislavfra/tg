import {RequestParamsType} from './types';
import {getSchemas} from './_utils/get-schemas';


const ENDPOINT = '/adder';

export const makeRequestConfig = ({additionalFetchParams,bodyParams}: RequestParamsType)=>({endpoint: ENDPOINT,responseSchema: getSchemas(bodyParams),body: bodyParams,isBatchRequest: true ,...additionalFetchParams});