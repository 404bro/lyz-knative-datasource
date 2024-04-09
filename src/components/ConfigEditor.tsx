import React, { ChangeEvent } from 'react';
import { InlineField, Input } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { MyDataSourceOptions } from '../types';

interface Props extends DataSourcePluginOptionsEditorProps<MyDataSourceOptions> { }

export function ConfigEditor(props: Props) {
  const { onOptionsChange, options } = props;
  const onAgentURLChange = (event: ChangeEvent<HTMLInputElement>) => {
    const jsonData = {
      ...options.jsonData,
      agentURL: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };

  const { jsonData } = options;

  return (
    <div className="gf-form-group">
      <InlineField label="Agent URL" labelWidth={24} >
        <Input
          onChange={onAgentURLChange}
          value={jsonData.agentURL || ''}
          width={50}
          placeholder='http://knative-agent.default.svc.cluster.local:9091'
        />
      </InlineField>
    </div>
  );
}
