export class Requirement {
  text: string;
  satisfied: boolean;

  constructor(text: string, satisfied: boolean) {
    this.text = text;
    this.satisfied = satisfied;
  }
}

export class Project {
  name: string;
  requirements: Requirement[];

  constructor(project: string, requirements: Requirement[]) {
    this.name = project;
    this.requirements = new Array<Requirement>();

    for (const elt of requirements) {
      this.requirements.push(elt);
    }
  }
}