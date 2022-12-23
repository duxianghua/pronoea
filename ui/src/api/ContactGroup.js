import request from '@/utils/request'

export function ListContactGroup(params) {
  return request({
    url: '/contactgroup/',
    method: 'get',
    params
  })
}

export function GetContactGroup(namespace, name) {
  return request({
    url: '/contactgroup/' + name,
    method: 'get'
  })
}

export function CreateContactGroup(data) {
  return request({
    url: '/contactgroup/' + data.metadata.name,
    method: 'post',
    data
  })
}

export function UpdateContactGroup(data) {
  return request({
    url: '/contactgroup/' + data.metadata.name,
    method: 'put',
    data
  })
}

export function DeleteContactGroup(name) {
  return request({
    url: '/contactgroup/' + name,
    method: 'delete'
  })
}

