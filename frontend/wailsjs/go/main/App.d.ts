// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {services} from '../models';
import {common} from '../models';

export function GetGithubActionLog(arg1:number):Promise<services.Msg>;

export function ReadConfig():Promise<common.Config>;

export function SendGithubAction(arg1:string,arg2:Array<string>):Promise<services.Msg>;

export function WriteConfig(arg1:common.Config):Promise<string>;