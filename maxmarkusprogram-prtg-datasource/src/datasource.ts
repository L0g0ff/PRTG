import { 
  DataSourceInstanceSettings, 
  ScopedVars, 
  DataQueryRequest,
  DataQueryResponse,
  MetricFindValue,
} from '@grafana/data';
import { 
  DataSourceWithBackend, 
  getTemplateSrv,
} from '@grafana/runtime';
import { Observable, from, merge, throwError } from 'rxjs';
import { catchError } from 'rxjs/operators';
import { processAnnotationEvents } from './datasource/annotations';
import { createStreamingQueryObservable } from './datasource/streaming';
import {
  MyQuery,
  MyDataSourceOptions,
  PRTGGroupListResponse,
  PRTGDeviceListResponse,
  PRTGSensorListResponse,
  PRTGChannelListResponse,
  QueryType,
} from './types'

export class DataSource extends DataSourceWithBackend<MyQuery, MyDataSourceOptions> {

  constructor(instanceSettings: DataSourceInstanceSettings<MyDataSourceOptions>) {
    super(instanceSettings);
  }

  applyTemplateVariables(query: MyQuery, scopedVars: ScopedVars) {
    const templateSrv = getTemplateSrv();
    const replace = (value?: string) => templateSrv.replace(value || '', scopedVars);
    const channel = replace(query.channel);
    const channelArray = query.channelArray?.map((item) => replace(item)).filter(Boolean) || [];

    return {
      ...query,
      group: replace(query.group),
      groupId: replace(query.groupId),
      device: replace(query.device),
      deviceId: replace(query.deviceId),
      sensor: replace(query.sensor),
      sensorId: replace(query.sensorId),
      channel,
      channelArray,
      manualObjectId: replace(query.manualObjectId),
    }
  }

  filterQuery(query: MyQuery): boolean {
    if (query.queryType === QueryType.Metrics) {
      return Boolean(query.sensorId && (query.channel || query.channelArray?.length));
    }

    return true;
  }

  async metricFindQuery(query: string | { query?: string }): Promise<MetricFindValue[]> {
    const queryText = typeof query === 'string' ? query : query.query || '';
    const [kind, ...argParts] = queryText.split(':');
    const rawArg = argParts.join(':');
    const arg = getTemplateSrv().replace(rawArg).trim();

    switch (kind.trim()) {
      case 'groups': {
        const response = await this.getGroups();
        return response.groups
          .filter((item) => item.group)
          .map((item) => ({ text: item.group, value: item.group }));
      }
      case 'devices': {
        if (!arg) {
          return [];
        }
        const response = await this.getDevices(arg);
        return response.devices
          .filter((item) => item.device)
          .map((item) => ({ text: item.device, value: item.device }));
      }
      case 'sensors': {
        if (!arg) {
          return [];
        }
        const response = await this.getSensors(arg);
        return response.sensors
          .filter((item) => item.sensor && item.objid)
          .map((item) => ({ text: item.sensor, value: String(item.objid) }));
      }
      case 'channels': {
        if (!arg) {
          return [];
        }
        const response = await this.getChannels(arg);
        const channelData = response.values?.[0];
        if (!channelData) {
          return [];
        }

        return Object.keys(channelData)
          .filter((key) => key !== 'datetime')
          .map((key) => ({ text: key, value: key }));
      }
      default:
        return [];
    }
  }

  async getGroups(): Promise<PRTGGroupListResponse> {
    return this.getResource('groups')
  }

  async getDevices(group: string): Promise<PRTGDeviceListResponse> {
    if (!group) {
      throw new Error('group is required')
    }
    return this.getResource(`devices/${encodeURIComponent(group)}`)
  }

  async getSensors(device: string): Promise<PRTGSensorListResponse> {
    if (!device) {
      throw new Error('device is required');
    }
    return this.getResource(`sensors/${encodeURIComponent(device)}`);
  }

  async getChannels(sensorId: string): Promise<PRTGChannelListResponse> {
    if (!sensorId) {
      throw new Error('sensorId is required');
    }
    return this.getResource(`channels/${encodeURIComponent(sensorId)}`);
  }

  annotations = {
    QueryEditor: undefined,
    processEvents: processAnnotationEvents,
  };

  query(request: DataQueryRequest<MyQuery>): Observable<DataQueryResponse> {
    // Only handle streaming for metrics queries
    const streamingTargets = request.targets.filter(
      query => query.isStreaming && query.queryType === QueryType.Metrics
    );
    const regularTargets = request.targets.filter(
      query => !query.isStreaming || query.queryType !== QueryType.Metrics
    );
    
    const observables: Array<Observable<DataQueryResponse>> = [];

    // Process streaming targets
    if (streamingTargets.length > 0) {
      streamingTargets.forEach((query) => {
        observables.push(createStreamingQueryObservable(this.uid, request, query));
      });
    }

    // Process regular targets
    if (regularTargets.length > 0) {
      observables.push(
        super.query({
          ...request,
          targets: regularTargets,
        }).pipe(
          catchError((err) => {
            console.error('Query error:', err);
            return throwError(() => err);
          })
        )
      );
    }

    // Return combined observables or empty data
    if (observables.length === 0) {
      return from([{ data: [] }]);
    }
    
    return merge(...observables);
  }

  // Stream control methods
  async getStreamStatus(streamId: string): Promise<any> {
    return this.getResource(`stream-status/${streamId}`);
  }

  async stopStream(streamId: string): Promise<void> {
    return this.getResource(`stop-stream/${streamId}`);
  }
}
