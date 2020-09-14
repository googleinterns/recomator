export class Requirement {
  text: string;
  satisfied: boolean;
  errorMessage: string;

  constructor(text: string, satisfied: boolean, errorMessage: string) {
    this.text = text;
    this.satisfied = satisfied;
    this.errorMessage = errorMessage;
  }
}

export class ProjectRequirement {
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
