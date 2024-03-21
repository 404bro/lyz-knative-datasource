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
  k8sUrl?: string;
  promUrl?: string;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface MySecureJsonData {
  k8sToken?: string;
}

export enum QueryTypeValue {
  OVERVIEW = 'overview',
}

export const QueryType: Array<SelectableValue<QueryTypeValue>> = [
  {
    label: 'overview',
    description: 'Overview Node Graph',
    value: QueryTypeValue.OVERVIEW
  },
]
