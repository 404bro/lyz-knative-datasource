import React, { ChangeEvent } from 'react';
import { InlineField, Input, Select } from '@grafana/ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { DataSource } from '../datasource';
import { MyDataSourceOptions, MyQuery, QueryType } from '../types';

type Props = QueryEditorProps<DataSource, MyQuery, MyDataSourceOptions>;

export function QueryEditor({ query, onChange, onRunQuery }: Props) {
  const onTypeChange = (value: SelectableValue<string>) => {
    onChange({ ...query, type: value.value || '' })
  };

  const onQueryTextChange = (event: ChangeEvent<HTMLInputElement>) => {
    onChange({ ...query, queryText: event.target.value });
  };

  const onConstantChange = (event: ChangeEvent<HTMLInputElement>) => {
    onChange({ ...query, constant: parseFloat(event.target.value) });
    // executes the query
    onRunQuery();
  };

  const { type, queryText, constant } = query;

  return (
    <div className="gf-form">
      <InlineField label="Type">
        <Select onChange={onTypeChange} value={type} width={24} options={QueryType} />
      </InlineField>
      {type !== "Service Map" && (
        <>
          <InlineField label="Constant">
            <Input onChange={onConstantChange} value={constant} width={8} type="number" step="0.1" />
          </InlineField>
          <InlineField label="Query Text" labelWidth={16} tooltip="Not used yet">
            <Input onChange={onQueryTextChange} value={queryText || ''} />
          </InlineField>
        </>
      )}
    </div>
  );
}
