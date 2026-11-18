import App from "./App";
import ReactDOMClient from 'react-dom/client';
import { TMfeProps } from '@milestone/one-shell-mfe-contract';
import "@milestone/ui-assets/fonts.css";

declare const __STANDALONE__: boolean;

if (__STANDALONE__) {
  const standaloneMockedProps: TMfeProps<{}> = {
      domElement: document.getElementById('root') as HTMLElement, //DOM Element used to mount this MFE into.
      matchRoute: (pathname: string) => {
          console.log(`matchRoute called with: ${pathname}`);
          return null;
      },
      navigate: (routeOrPath: string) => { // Shell SDK navigate function
          console.log(`navigate called with: ${routeOrPath}`);
      },
      mountParcel: (async (module: any, props: any, mfeDef: any) => { // Shell SDK mountParcel function, used to mount MFEs withing a container.
          console.log('mountParcel called', { module, props, mfeDef });
          return {} as any;
      }) as TMfeProps['mountParcel'],
      loadRemoteModule: async (mfeDef, container) => { // Shell SDK loadRemoteModule function used to load remote MFE modules
          console.log('loadRemoteModule called', { mfeDef, container });
          return {} as any;
      },
      findRemoteMfeDefinition: (predicate) => { // Shell SDK findRemoteMfeDefinition function used to find MFE definitions based on a predicate function
          console.log('findRemoteMfeDefinition called', predicate);
          return undefined;
      },
      mountPlaceholderMfes: ((placeholdersConfig) => { // Shell SDK function that mounts every MFE registered for each `position` into the placeholder DOM elements declared by this MFE.
        console.log('mountPlaceholderMfes called', placeholdersConfig);
        return new Map();
      }) as TMfeProps['mountPlaceholderMfes'],
      localMfeConfig: [], // Local MFE configuration
      localMfeModuleRegistry: {} as TMfeProps['localMfeModuleRegistry'], // Local MFE module registry
      executeAction: async (event, inWorkflowContext, urlParams) => { // Shell SDK executeAction function used to run an action provided by the native host (relayed to the host, on the workflow's channel inside a workflow)
          console.log('executeAction called', { event, inWorkflowContext, urlParams });
          return {}; // No shell, so no workflow: resolves like an action opened outside a workflow
      },
  }
  const rootElement = document.getElementById('root');
  if (rootElement) {
    const root = ReactDOMClient.createRoot(rootElement);
    root.render(<App {...standaloneMockedProps} />);
  }
}

export default App;