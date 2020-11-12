import {SCHEMAS} from '../../_schemas';
import {BatchedParamsType} from '../types';


export const getSchemas = (options: Array<BatchedParamsType>)=>options.map(({
method
})=>SCHEMAS[method]);