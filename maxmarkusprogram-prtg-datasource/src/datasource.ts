import { 
  DataSourceInstanceSettings, 
  ScopedVars, 
  AnnotationEvent,
  DataFrame,
  DataQueryRequest,
  DataQueryResponse,
  LiveChannelScope,
  MetricFindValue,
} from '@grafana/data';
import { 
  DataSourceWithBackend, 
  getTemplateSrv,
  getGrafanaLiveSrv,
} from '@grafana/runtime';
import { Observable, from, merge, throwError } from 'rxjs';
import { catchError, map } from 'rxjs/operators';
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
    processEvents: (anno: any, data: DataFrame[]): Observable<AnnotationEvent[]> => {
      const events: AnnotationEvent[] = [];
      
      data.forEach((frame) => {
        const timeField = frame.fields.find((field) => field.name === 'Time');
        const valueField = frame.fields.find((field) => field.name === 'Value');
        
        if (timeField && valueField) {
          const firstTime = timeField.values[0];
          const lastTime = timeField.values[timeField.values.length - 1];
          const firstValue = valueField.values[0];
          const panelId = typeof anno.panelId === 'number' ? anno.panelId : undefined;
          const source = frame.name || 'PRTG Channel';

          events.push({
            time: firstTime,
            timeEnd: lastTime !== firstTime ? lastTime : undefined,
            title: source,
            text: `Value: ${firstValue}`,
            tags: ['prtg', `value:${firstValue}`, `source:${source}`],
            panelId: panelId
          });
        }
      });

      return from([events]);
    },
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
        // Add panelId to query for stream ID generation
        const queryWithPanelId = {
          ...query,
          panelId: request.panelId?.toString()
        };
        
        // Create a unique, stable stream ID
        const streamId = this.getStreamId(queryWithPanelId);
        const streamPath = `prtg-stream/${streamId}`;
        
        // Set up the data stream
        const streamObs = getGrafanaLiveSrv().getDataStream({
          addr: {
            scope: LiveChannelScope.DataSource,
            namespace: this.uid,
            path: streamPath,
            data: {
              ...query,
              streamId,
              panelId: request.panelId?.toString(),
              queryId: query.refId,
              timeRange: {
                from: request.range.from.valueOf(),
                to: request.range.to.valueOf(),
              },
              // Use provided values or defaults
              cacheTime: query.cacheTime,
              updateMode: query.updateMode,
              bufferSize: query.bufferSize,
            },
          },
        }).pipe(
          map((response) => {
            // Enhance frame with streaming metadata
            const frameData = response.data || [];
            frameData.forEach((frame) => {
              if (frame && frame.meta) {
                frame.meta = {
                  ...frame.meta,
                  streaming: true,
                  streamId,
                  preferredVisualisationType: 'graph',
                };
              }
            });
            return { data: frameData };
          }),
          catchError((err) => {
            console.error('Stream error:', err);
            return throwError(() => new Error(`Streaming error: ${err.message || 'Unknown error'}`));
          })
        );
        
        observables.push(streamObs);
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

  // Improved stream ID generation for better stability
  private getStreamId(query: MyQuery & { panelId?: string }): string {
    const components = [
      query.panelId || 'default',
      query.refId || 'A',
      query.sensorId || '',
      Array.isArray(query.channelArray) && query.channelArray.length > 0 
        ? query.channelArray.join('-') 
        : query.channel || '',
    ];
    return components.filter(Boolean).join('_');
  }

  // Stream control methods
  async getStreamStatus(streamId: string): Promise<any> {
    return this.getResource(`stream-status/${streamId}`);
  }

  async stopStream(streamId: string): Promise<void> {
    return this.getResource(`stop-stream/${streamId}`);
  }
}
