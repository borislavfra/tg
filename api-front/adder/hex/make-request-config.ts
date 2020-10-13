import {responseSchema} from './response-schema';
import {RequestParamsType} from './types';

const ENDPOINT = '/adder/hex'

export const makeRequestConfig = ({additionalFetchParams,bodyParams}: RequestParamsType)=>({endpoint: ENDPOINT,responseSchema,body: {params: bodyParams},...additionalFetchParams});