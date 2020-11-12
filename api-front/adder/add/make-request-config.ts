import {SCHEMAS} from '../_schemas';
import {RequestParamsType} from './types';

const ENDPOINT = '/adder/add'

export const makeRequestConfig = ({additionalFetchParams,bodyParams}: RequestParamsType)=>({endpoint: ENDPOINT,responseSchema: SCHEMAS.add,body: {params: bodyParams},...additionalFetchParams});