import React, { useEffect, useState, useMemo, useCallback, ChangeEvent, useRef } from 'react';
import {
  InlineField,
  Combobox,
  Stack,
  AsyncMultiSelect,
} from '@grafana/ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data'
import type { ComboboxOption } from '@grafana/ui';
import { DataSource } from '../datasource'
import {
  MyDataSourceOptions, MyQuery, queryTypeOptions, QueryType, propertyList, filterPropertyList, manualApiMethods
} from '../types'
import { DisplayOptions } from './query-editor/DisplayOptions';
import { ManualQueryOptions } from './query-editor/ManualQueryOptions';
import { PropertyOptions } from './query-editor/PropertyOptions';
import { StreamingOptions } from './query-editor/StreamingOptions';

type Props = QueryEditorProps<DataSource, MyQuery, MyDataSourceOptions>

export function QueryEditor({ query, onChange, onRunQuery, datasource }: Props) {
  const prevQueryRef = useRef<MyQuery | null>(null);
  const runQueryIfChanged = useCallback(() => {
    const currentQuery = JSON.stringify({ ...query, refId: query.refId }); // Include refId in comparison
    const prevQuery = JSON.stringify(prevQueryRef.current);

    if (currentQuery !== prevQuery) {
      prevQueryRef.current = query;
      onRunQuery();
    }
  }, [query, onRunQuery]);

  const isMetricsMode = query.queryType === QueryType.Metrics
  const isRawMode = query.queryType === QueryType.Raw
  const isTextMode = query.queryType === QueryType.Text
  const isManualMode = query.queryType === QueryType.Manual

  /* ===================================================== HOOKS ============================================================*/
  const [group, setGroup] = useState<string>(query.group || '')
  const [device, setDevice] = useState<string>(query.device || '')
  const [sensor, setSensor] = useState<string>(query.sensor || '')
  //@ts-ignore
  const [channel, setChannel] = useState<string>(query.channel || '')
  const [channelQuery, setChannelQuery] = useState<string[]>(query.channelArray || [])
  const [sensorId, setSensorId] = useState<string>(query.sensorId || '')
  const [manualMethod, setManualMethod] = useState<string>(query.manualMethod || '');
  const [manualObjectId, setManualObjectId] = useState<string>(query.manualObjectId || '');
  const [streamIntervalValue, setStreamIntervalValue] = useState<string>(String(query.streamInterval || 2500));

  const [lists, setLists] = useState({
    groups: [] as Array<ComboboxOption<string>>,
    devices: [] as Array<ComboboxOption<string>>,
    sensors: [] as Array<ComboboxOption<string>>,
    channels: [] as Array<SelectableValue<string>>,
    values: [] as Array<SelectableValue<string>>,
    properties: [] as Array<SelectableValue<string>>,
    filterProperties: [] as Array<SelectableValue<string>>,
  })

  const [isLoading, setIsLoading] = useState(false)


  /* ================================================== SORT ================================================== */
  lists.groups.sort((a, b) => (a.label ?? '').localeCompare(b.label ?? ''))
  lists.devices.sort((a, b) => (a.label ?? '').localeCompare(b.label ?? ''))
  lists.sensors.sort((a, b) => (a.label ?? '').localeCompare(b.label ?? ''))
  lists.channels.sort((a, b) => (a.label ?? '').localeCompare(b.label ?? ''))


  /* ================================================== FETCH GROUPS ================================================== */
  useEffect(() => {
    async function fetchGroups() {
      setIsLoading(true)
      try {
        const response = await datasource.getGroups()
        if (response && Array.isArray(response.groups)) {
          const groupOptions = response.groups.map((group) => ({
            label: group.group,
            value: group.group.toString(),
          }))

          // Batch state updates to avoid ACT warnings in tests
          setTimeout(() => {
            setLists((prev) => ({
              ...prev,
              groups: groupOptions,
            }))
            setIsLoading(false)
          }, 0)
        } else {
          console.error('Invalid response format:', response)
          setTimeout(() => {
            setLists((prev) => ({
              ...prev,
              groups: [],
            }))
            setIsLoading(false)
          }, 0)
        }
      } catch (error) {
        console.error('Error fetching groups:', error)
        setTimeout(() => {
          setLists((prev) => ({
            ...prev,
            groups: [],
          }))
          setIsLoading(false)
        }, 0)
      }
    }
    fetchGroups()
  }, [datasource])
  /* ================================================== FETCH DEVICES ================================================== */
  useEffect(() => {
    async function fetchDevices() {
      if (!group) { return };

      setIsLoading(true)
      try {
        const response = await datasource.getDevices(group)
        if (response && Array.isArray(response.devices)) {
          const filteredDevices = group ? response.devices.filter((device) => device.group === group) : response.devices

          const deviceOptions = filteredDevices.map((device) => ({
            label: device.device,
            value: device.device.toString(),
          }))
          setLists((prev) => ({
            ...prev,
            devices: deviceOptions,
          }))
        } else {
          console.error('Invalid devices response format:', response)
          setLists((prev) => ({
            ...prev,
            devices: [],
          }))
        }
      } catch (error) {
        console.error('Error fetching devices:', error)
        setLists((prev) => ({
          ...prev,
          devices: [],
        }))
      } finally {
        setIsLoading(false)
      }
    }
    fetchDevices()
  }, [datasource, group])
  /* ================================================== FETCH SENSOR ================================================== */
  useEffect(() => {
    async function fetchSensors() {
      if (!device) { return };

      setIsLoading(true)
      try {
        const response = await datasource.getSensors(device)
        if (response && Array.isArray(response.sensors)) {
          const filteredSensors = device
            ? response.sensors.filter((sensor) => sensor.device === device)
            : response.sensors
          const sensorOptions = filteredSensors.map((sensor) => ({
            label: sensor.sensor,
            value: sensor.sensor.toString(),
          }))
          setLists((prev) => ({
            ...prev,
            sensors: sensorOptions,
          }))
        } else {
          console.error('Invalid sensors response format:', response)
          setLists((prev) => ({
            ...prev,
            sensors: [],
          }))
        }
      } catch (error) {
        console.error('Error fetching sensors:', error)
        setLists((prev) => ({
          ...prev,
          sensors: [],
        }))
      } finally {
        setIsLoading(false)
      }
    }
    fetchSensors()
  }, [datasource, device])
  /* ==================================================  FETCH CHANNEL ==================================================   */
  useEffect(() => {
    async function fetchChannels() {
      if (!sensorId) {
        setLists((prev) => ({
          ...prev,
          channels: [],
        }));
        return;
      }

      setIsLoading(true);
      try {
        const response = await datasource.getChannels(sensorId);
        if (!response) {
          console.error('Empty response received');
          setLists((prev) => ({
            ...prev,
            channels: [],
          }));
          return;
        }

        if (response.values && Array.isArray(response.values) && response.values.length > 0) {
          const channelData = response.values[0] || {};

          const channelOptions = Object.entries(channelData)
            .filter(([key]) => key !== 'datetime')
            .map(([key]) => ({
              label: key,
              value: key,
            }));

          setLists((prev) => ({
            ...prev,
            channels: channelOptions,
          }));

          if (query.channel && channelOptions.some(opt => opt.value === query.channel)) {
            setChannel(query.channel);
          }
        } else {
          console.warn('No channel data found in response');
          setLists((prev) => ({
            ...prev,
            channels: [],
          }));
        }

      } catch (error) {
        console.error('Error fetching channels:', error);
        setLists((prev) => ({
          ...prev,
          channels: [],
        }));
      } finally {
        setIsLoading(false);
      }
    }

    fetchChannels();
  }, [datasource, sensorId, query.channel]);

  useEffect(() => {
    if (isTextMode || isRawMode) {
      const propertyOptions = propertyList.map((item) => ({
        label: item.visible_name,
        value: item.name,
      }));

      // Filter property options
      const filterPropertyOptions = filterPropertyList.map((item) => ({
        label: item.visible_name,
        value: item.name,
      }));

      setLists((prev) => ({
        ...prev,
        properties: propertyOptions,
        filterProperties: filterPropertyOptions,
      }));
    }
  }, [isTextMode, isRawMode]);
  /* ==================================================  INITIAL VALUES  ================================================== */
  useEffect(() => {
    setGroup((prev) => query.group ?? prev);
    setDevice((prev) => query.device ?? prev);
    setSensor((prev) => query.sensor ?? prev);
    setChannel((prev) => query.channel ?? prev);
    setSensorId((prev) => query.sensorId ?? prev);
    setManualMethod((prev) => query.manualMethod ?? prev);
    setManualObjectId((prev) => query.manualObjectId ?? prev);
    setStreamIntervalValue(String(query.streamInterval || 2500));
    // Add this line to restore channel selections
    setChannelQuery((prev) => query.channelArray || prev || []);
  }, [query]);


  /* ==================================================  FIND IDs ================================================= */
  const findGroupId = useCallback(async (groupName: string) => {
    try {
      const response = await datasource.getGroups()
      if (response && Array.isArray(response.groups)) {
        const group = response.groups.find((g) => g.group === groupName)
        if (group) {
          return group.objid.toString()
        }
      }
    } catch (error) {
      console.error('Error finding group ID:', error)
    }
    return ''
  }, [datasource])

  const findDeviceId = useCallback(async (deviceName: string) => {
    try {
      const response = await datasource.getDevices(group)
      if (response && Array.isArray(response.devices)) {
        const device = response.devices.find((d) => d.device === deviceName)
        if (device) {
          return device.objid.toString()
        }
      }
    } catch (error) {
      console.error('Error finding device ID:', error)
    }
    return ''
  }, [datasource, group])

  const findSensorObjid = useCallback(async (sensorName: string) => {
    try {
      const response = await datasource.getSensors(device)
      if (response && Array.isArray(response.sensors)) {
        const sensor = response.sensors.find((s) => s.sensor === sensorName)
        if (sensor) {
          setSensorId(sensor.objid.toString())
          return sensor.objid.toString()
        } else {
          console.error('Sensor not found:', sensorName)
        }
      } else {
        console.error('Invalid response format:', response)
      }
    } catch (error) {
      console.error('Error fetching sensors:', error)
    }
    return ''
  }, [datasource, device, setSensorId])
  /* ==================================================  USE MEMO  ==================================================  */

  const groupOptions = useMemo(() => lists.groups, [lists.groups]);
  const deviceOptions = useMemo(() => lists.devices, [lists.devices]);
  const sensorOptions = useMemo(() => lists.sensors, [lists.sensors]);

  // Add new memoized selected values
  const selectedGroup = useMemo(() => {
    return groupOptions.find(option => option.value === group) || (group ? { label: group, value: group } : null);
  }, [groupOptions, group]);

  const selectedDevice = useMemo(() => {
    return deviceOptions.find(option => option.value === device) || (device ? { label: device, value: device } : null);
  }, [deviceOptions, device]);

  const selectedSensor = useMemo(() => {
    return sensorOptions.find(option => option.value === sensor) || (sensor ? { label: sensor, value: sensor } : null);
  }, [sensorOptions, sensor]);

  // Add new loadChannelOptions function with useMemo
  const loadChannelOptions = useMemo(() => async () => {
    if (!sensorId) {
      return [];
    }

    try {
      const response = await datasource.getChannels(sensorId);

      if (!response) {
        console.warn('No response received from getChannels');
        return [];
      }

      // Check if response has the expected structure
      if (typeof response === 'object' && 'values' in response) {
        const values = response.values;
        if (!Array.isArray(values) || values.length === 0) {
          console.warn('No channel values found in response');
          return [];
        }

        const channelData = values[0];
        if (typeof channelData !== 'object') {
          console.warn('Invalid channel data format');
          return [];
        }

        return Object.keys(channelData)
          .filter(key => key !== 'datetime')
          .map(key => ({
            label: key,
            value: key,
          }));
      }

      console.warn('Unexpected response format:', response);
      return [];
    } catch (error: any) {
      console.error('Error loading channels:', error?.message || error);
      return [];
    }
  }, [sensorId, datasource]);
  /* ==================================================  EVENT HANDLERS ==================================================  */

  /* ==================================================  QUERY  ==================================================  */
  /* ==================================================  ONQUERYTYPESCHANGE ==================================================  */
  const onQueryTypeChange = useCallback((option: ComboboxOption<string> | null) => {
    if (option?.value) {
      onChange({ ...query, queryType: option.value as QueryType });
      runQueryIfChanged();
    }
  }, [query, onChange, runQueryIfChanged]);

  /* ==================================================  ONGROUPCHANGE ==================================================  */
  const onGroupChange = useCallback(async (option: ComboboxOption<string> | null) => {
    if (!option?.value) return;

    const groupObjId = await findGroupId(option.value);
    setGroup(option.value);

    const updatedQuery = {
      ...query,
      group: option.value,
      groupId: groupObjId,
    };
    onChange(updatedQuery);
    setLists(prev => ({ ...prev, devices: [], sensors: [], channels: [] }));
    runQueryIfChanged();
  }, [query, onChange, runQueryIfChanged, findGroupId]);

  /* ==================================================  ONDEVICECHANGE ================================================= */
  const onDeviceChange = useCallback(async (option: ComboboxOption<string> | null) => {
    if (!option?.value) return;

    const deviceObjId = await findDeviceId(option.value);

    setDevice(option.value);
    const updatedQuery = {
      ...query,
      device: option.value,
      deviceId: deviceObjId,
    };
    onChange(updatedQuery);
    setLists(prev => ({ ...prev, sensors: [], channels: [] }));
    runQueryIfChanged();
  }, [query, onChange, runQueryIfChanged, findDeviceId]);
  /* ==================================================  ONSENSORCHANGE ==================================================  */
  const onSensorChange = useCallback(async (option: ComboboxOption<string> | null) => {
    if (!option?.value) {
      return;
    }

    const sensorObjId = await findSensorObjid(option.value);

    setSensor(option.value);
    setSensorId(sensorObjId);
    setLists(prev => ({ ...prev, channels: [] }));

    const updatedQuery = {
      ...query,
      sensor: option.value,
      sensorId: sensorObjId,
    };
    onChange(updatedQuery);

    runQueryIfChanged();
  }, [query, onChange, runQueryIfChanged, findSensorObjid]);  /* ==================================================  ONCHANNELCHANGE ==================================================  */
  const onChannelChange = useCallback((values: Array<SelectableValue<string>>) => {
    const selectedChannels = values.map(v => v.value!);

    // Update local state
    setChannelQuery(selectedChannels);

    // CRITICAL: Update query to include ALL selected channels in a SINGLE query
    // This prevents Grafana from creating multiple queries (refId A, B, C...)
    const updatedQuery = {
      ...query,
      channel: selectedChannels[0] || '', // First channel for backward compatibility
      channelArray: selectedChannels, // ALL selected channels in one array
      // Generate series names for each channel
      seriesNames: selectedChannels.map(channel =>
        `${query.sensor || 'Sensor'} - ${channel}`
      ),
    };

    onChange(updatedQuery);    // Only trigger query execution if we have channels selected
    // Don't use runQueryIfChanged() as it might create duplicate queries
    if (selectedChannels.length > 0) {
      // Use timeout to ensure state is updated before running query
      setTimeout(() => {
        onRunQuery();
      }, 0);
    }
  }, [query, onChange, onRunQuery]);
  /* ==================================================  ON INCLUDE GROUP NAME ==================================================  */
  const onIncludeGroupName = (event: ChangeEvent<HTMLInputElement>) => {
    const updatedQuery = { ...query, includeGroupName: event.currentTarget.checked };
    onChange(updatedQuery);
    runQueryIfChanged();
  }

  /* ==================================================  ON INCLUDE DEVICE NAME ==================================================  */
  const onIncludeDeviceName = (event: React.ChangeEvent<HTMLInputElement>) => {
    const updatedQuery = { ...query, includeDeviceName: event.currentTarget.checked };
    onChange(updatedQuery);
    runQueryIfChanged();
  }

  /* ==================================================  ON INCLUDE SENSOR NAME ==================================================  */
  const onIncludeSensorName = (event: ChangeEvent<HTMLInputElement>) => {
    const updatedQuery = { ...query, includeSensorName: event.currentTarget.checked };
    onChange(updatedQuery);
    runQueryIfChanged();
  }  /* ==================================================  ON MANUAL OBJECT ID CHANGE ==================================================  */
  const onManualObjectIdChange = (event: ChangeEvent<HTMLInputElement>) => {
    const value = event.currentTarget.value;
    setManualObjectId(value);
    const updatedQuery = {
      ...query,
      manualObjectId: value,
    };
    onChange(updatedQuery);
    runQueryIfChanged();
  };

  const onManualMethodChange = useCallback((value: string) => {
    setManualMethod(value);
    const updatedQuery = {
      ...query,
      manualMethod: value,
    };
    onChange(updatedQuery);
    runQueryIfChanged();
  }, [query, onChange, runQueryIfChanged]);

  const onPropertyChange = useCallback((value: string) => {
    const updatedQuery = { ...query, property: value };
    onChange(updatedQuery);
    runQueryIfChanged();
  }, [query, onChange, runQueryIfChanged]);

  const onFilterPropertyChange = useCallback((value: string) => {
    const updatedQuery = { ...query, filterProperty: value };
    onChange(updatedQuery);
    runQueryIfChanged();
  }, [query, onChange, runQueryIfChanged]);

  /* ==================================================  STREAM INTERVAL HANDLERS ==================================================  */
  const handleStreamIntervalChange = useCallback((e: ChangeEvent<HTMLInputElement>) => {
    const value = e.currentTarget.value;
    setStreamIntervalValue(value);
  }, []);

  const handleStreamIntervalBlur = useCallback(() => {
    const interval = Math.max(0, Math.min(60000, parseInt(streamIntervalValue, 10) || 2500));
    const updatedQuery = {
      ...query,
      streamInterval: interval,
    };
    onChange(updatedQuery);
    runQueryIfChanged();
  }, [streamIntervalValue, query, onChange, runQueryIfChanged]);

  const onStreamingToggle = useCallback((event: ChangeEvent<HTMLInputElement>) => {
    const isStreaming = event.currentTarget.checked;
    const streamInterval = isStreaming ? (query.streamInterval || 2500) : undefined;
    const updatedQuery = {
      ...query,
      isStreaming,
      streamInterval,
    };
    onChange(updatedQuery);
    runQueryIfChanged();
  }, [query, onChange, runQueryIfChanged]);



  /* ================================================== DESTRUCTURING ================================================== */
  // Set default streaming values
  useEffect(() => {
    if (query.isStreaming === undefined) {
      const updatedQuery = {
        ...query,
        isStreaming: false,
        streamInterval: 2500, // Default interval 5ms (2,5 seconds)
      };
      onChange(updatedQuery);
    }
  }, [query, onChange]);

  /* ================================================== RENDER ================================================== */
  return (
    <Stack direction="column" gap={2}>
      <Stack direction="row" gap={2}>
        <Stack direction="column" gap={1}>          <InlineField label="Query Type" labelWidth={20} grow>
          <Combobox
            id='query-editor-queryType'
            options={queryTypeOptions}
            value={query.queryType}
            onChange={onQueryTypeChange}
            width={47}
          />
        </InlineField>

          <InlineField label="Group" labelWidth={20} grow>
            <Combobox
              id='query-editor-group'
              loading={isLoading}
              options={groupOptions}
              value={selectedGroup}
              onChange={onGroupChange}
              width={47}
              createCustomValue={true}
              isClearable={true}
              invalid={!query.queryType}
              placeholder="Select Group or type '*'"
            />
          </InlineField>

          <InlineField label="Device" labelWidth={20} grow>
            <Combobox
              id='query-editor-device'
              loading={!lists.devices.length && !!query.group}
              options={deviceOptions}
              value={selectedDevice}
              onChange={onDeviceChange}
              width={47}
              createCustomValue={true}
              isClearable={true}
              invalid={!query.group}
              placeholder="Select Device or type '*'"
            />
          </InlineField>
        </Stack>

        <Stack direction="column" gap={1}>          <InlineField label="Sensor" labelWidth={20} grow>
          <Combobox
            id='query-editor-sensor'
            loading={!lists.sensors.length && !!query.device}
            options={sensorOptions}
            value={selectedSensor}
            onChange={onSensorChange}
            width={47}
            createCustomValue={true}
            isClearable={true}
            invalid={!query.device}
            placeholder="Select Sensor or type '*'"
          />
        </InlineField><InlineField label="Channel" labelWidth={20} grow>
            <AsyncMultiSelect
              id='query-editor-channel'
              key={sensorId}
              loadOptions={loadChannelOptions}
              defaultOptions={true}
              value={channelQuery.map(c => ({
                label: c,
                value: c,
              }))}
              onChange={onChannelChange}
              width={47}
              placeholder={sensorId ? "Select Channels (multiple allowed)" : "First select a sensor"}
              isClearable
              isDisabled={!sensorId}
              noOptionsMessage="No channels available"
            />
          </InlineField>
        </Stack>
      </Stack>


      {/* Show display name options for both Metrics and Streaming */}
      {(isMetricsMode || query.isStreaming || isRawMode || isTextMode) && (
        <DisplayOptions
          query={query}
          onIncludeGroupName={onIncludeGroupName}
          onIncludeDeviceName={onIncludeDeviceName}
          onIncludeSensorName={onIncludeSensorName}
        />
      )}
      
      {/* Options for Text and Raw modes */}
      {(isTextMode || isRawMode) && (
        <PropertyOptions
          property={query.property}
          filterProperty={query.filterProperty}
          properties={lists.properties}
          filterProperties={lists.filterProperties}
          onPropertyChange={onPropertyChange}
          onFilterPropertyChange={onFilterPropertyChange}
        />
      )}      {/* Manual API Query Section */}
      {isManualMode && (
        <ManualQueryOptions
          manualMethods={manualApiMethods}
          manualMethod={manualMethod}
          manualObjectId={manualObjectId}
          sensorId={sensorId}
          onManualMethodChange={onManualMethodChange}
          onManualObjectIdChange={onManualObjectIdChange}
        />
      )}

      {/* Always show streaming options */}
      <StreamingOptions
        query={query}
        streamIntervalValue={streamIntervalValue}
        onStreamingToggle={onStreamingToggle}
        onStreamIntervalChange={handleStreamIntervalChange}
        onStreamIntervalBlur={handleStreamIntervalBlur}
      />

    </Stack>
  )
}
