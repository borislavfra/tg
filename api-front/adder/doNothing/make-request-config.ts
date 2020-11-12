import {SCHEMAS} from '../_schemas';
import {RequestParamsType} from './types';

const ENDPOINT = '/adder/doNothing'

export const makeRequestConfig = ({additionalFetchParams,bodyParams}: RequestParamsType)=>({endpoint: ENDPOINT,responseSchema: SCHEMAS.doNothing,body: {params: bodyParams},...additionalFetchParams});