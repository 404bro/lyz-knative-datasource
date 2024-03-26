import React, { ChangeEvent } from 'react';
import { InlineField, Input, SecretInput } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { MyDataSourceOptions, MySecureJsonData } from '../types';

interface Props extends DataSourcePluginOptionsEditorProps<MyDataSourceOptions> { }

export function ConfigEditor(props: Props) {
  const { onOptionsChange, options } = props;
  const onK8sUrlChange = (event: ChangeEvent<HTMLInputElement>) => {
    const jsonData = {
      ...options.jsonData,
      k8sUrl: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };
  const onPromUrlChange = (event: ChangeEvent<HTMLInputElement>) => {
    const jsonData = {
      ...options.jsonData,
      promUrl: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };
  const onJaegerUrlChange = (event: ChangeEvent<HTMLInputElement>) => {
    const jsonData = {
      ...options.jsonData,
      jaegerUrl: event.target.value,
    };
    onOptionsChange({ ...options, jsonData })
  };

  // Secure field (only sent to the backend)
  const onK8sTokenChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      secureJsonData: {
        k8sToken: event.target.value,
      },
    });
  };

  const onResetK8sToken = () => {
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        k8sToken: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        k8sToken: '',
      },
    });
  };

  const { jsonData, secureJsonFields } = options;
  const secureJsonData = (options.secureJsonData || {}) as MySecureJsonData;

  return (
    <div className="gf-form-group">
      <InlineField label="Kubernetes API Server URL" labelWidth={24}>
        <Input
          onChange={onK8sUrlChange}
          value={jsonData.k8sUrl || ''}
          width={50}
          placeholder='https://localhost:6443'
        />
      </InlineField>
      <InlineField label="Kubernetes API Server Token" labelWidth={24}>
        <SecretInput
          isConfigured={(secureJsonFields && secureJsonFields.k8sToken) as boolean}
          value={secureJsonData.k8sToken || ''}
          width={50}
          onReset={onResetK8sToken}
          onChange={onK8sTokenChange}
        />
      </InlineField>
      <InlineField label="Prometheus URL" labelWidth={24}>
        <Input
          onChange={onPromUrlChange}
          value={jsonData.promUrl || ''}
          width={50}
          placeholder='http://localhost:9090'
        />
      </InlineField>
      <InlineField label="Jaeger URL" labelWidth={24} >
        <Input
          onChange={onJaegerUrlChange}
          value={jsonData.jaegerUrl || ''}
          width={50}
          placeholder='http://localhost:16686'
        />
      </InlineField>
    </div>
  );
}
