import { DataSourceJsonData, SelectableValue } from '@grafana/data';
import { DataQuery } from '@grafana/schema';

export interface MyQuery extends DataQuery {
  type: string
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
}

export const QueryType: Array<SelectableValue<QueryTypeValue>> = [
  {
    label: 'Service Map',
    description: 'Service Map (node graph)',
    value: QueryTypeValue.SERVICE_MAP
  },
]
