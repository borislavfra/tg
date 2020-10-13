import {request as adderAddRequest} from './adder/add';
import {request as adderHexRequest} from './adder/hex';


export const APICreator = {
adder: {
AddRequest: adderAddRequest,
HexRequest: adderHexRequest,
},
}