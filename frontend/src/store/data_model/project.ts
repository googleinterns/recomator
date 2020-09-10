export class Project {
  name: string;
  requirements: string[];

  constructor(project: string, requirements: string[]) {
    this.name = project;
    this.requirements = new Array<string>();

    for (const elt of requirements) {
      this.requirements.push(elt);
    }
  }
}