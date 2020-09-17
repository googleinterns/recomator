export interface IProjectsStoreState {
    projects: Project[];
    projectsSelected: Project[];
    loading: boolean;
    loaded: boolean;
  }
  
  export function projectsStoreStateFactory(): IProjectsStoreState {
    return {
      projects: [],
      projectsSelected: [],
      loading: false,
      loaded: false,
    };
  }