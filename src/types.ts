import { DataSourceJsonData, SelectableValue } from '@grafana/data';
import { DataQuery } from '@grafana/schema';

export interface MyQuery extends DataQuery {
  type: string
  name: string
  namespace: string
  queryText?: string;
  constant: number;
}

export const DEFAULT_QUERY: Partial<MyQuery> = {
  constant: 6.5,
};

/**
 * These are options configured for each DataSource instance
 */
export interface MyDataSourceOptions extends DataSourceJsonData {
  agentURL?: string;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface MySecureJsonData {
}

export enum QueryTypeValue {
  SERVICE_MAP = 'Service Map',
  SERVICES = 'Services',
  ANALYTICS_SERVICES = 'Analytics Services',
  ANALYTICS_TRACES = 'Analytics Traces',
  ANALYTICS_MAP = 'Analytics Map',
}

export const QueryType: Array<SelectableValue<QueryTypeValue>> = [
  {
    label: 'Service Map',
    description: 'Service Map (node graph)',
    value: QueryTypeValue.SERVICE_MAP
  },
  {
    label: 'Services',
    description: 'Services (table)',
    value: QueryTypeValue.SERVICES
  },
  {
    label: 'Analytics Services',
    description: 'Analytics Services (table)',
    value: QueryTypeValue.ANALYTICS_SERVICES
  },
  {
    label: 'Analytics Traces',
    description: 'Analytics Traces (table)',
    value: QueryTypeValue.ANALYTICS_TRACES
  },
  {
    label: 'Analytics Map',
    description: 'Analytics Map (node graph)',
    value: QueryTypeValue.ANALYTICS_MAP
  },
]
