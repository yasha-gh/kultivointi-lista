// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {context} from '../models';
import {main} from '../models';

export function GetContext():Promise<context.Context>;

export function GetDbAbsolutePath():Promise<string>;

export function GetSelf():Promise<main.Settings>;

export function OnSync():Promise<void>;

export function Save():Promise<void>;

export function SetContext(arg1:context.Context):Promise<void>;

export function Sync(arg1:main.Settings):Promise<void>;