import {JSONRPCRequest,IResponse} from '@mihanizm56/fetch-api';
import {makeRequestConfig} from './make-request-config';
import {RequestParamsType} from './types';

export const request = (values: RequestParamsType): Promise<IResponse>=>
new JSONRPCRequest().makeRequest(makeRequestConfig(values));