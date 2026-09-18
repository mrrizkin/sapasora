import account from './account';
import apikey from './apikey';
import auth from './auth';
import dashboard from './dashboard';
import device from './device';
import devicetoken from './devicetoken';
import gateway from './gateway';
import role from './role';

export const controllers = {
  account: Object.assign(account, account),
  apikey: Object.assign(apikey, apikey),
  auth: Object.assign(auth, auth),
  dashboard: Object.assign(dashboard, dashboard),
  device: Object.assign(device, device),
  devicetoken: Object.assign(devicetoken, devicetoken),
  gateway: Object.assign(gateway, gateway),
  role: Object.assign(role, role),
};

export default controllers;
