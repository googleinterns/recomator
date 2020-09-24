import { mount } from "@vue/test-utils";
import "@testing-library/jest-dom";
import Footer from "@/components/Footer.vue";
import { rootStoreFactory } from "@/store/root_store";
import { IRootStoreState } from "@/store/root_state";
import { RecommendationExtra } from "@/store/data_model/recommendation_extra";
import {
  freshSampleRawRecommendation,
  freshSavingRawRecommendation,
  freshPerformanceRawRecommendation
} from "./sample_recommendation";
import vuetify from "@/plugins/vuetify";
import { Store } from "vuex";

// We want all fetches to be mocked to do nothing here,
//  so that tests do not depenend on network requests
import { enableFetchMocks } from "jest-fetch-mock";
enableFetchMocks();

describe("Footer", () => {
  let recommendation: RecommendationExtra;
  let savingRecommendation: RecommendationExtra;
  let performanceRecommendation: RecommendationExtra;
  let fakeStore: Store<IRootStoreState>;

  beforeEach(() => {
    recommendation = new RecommendationExtra(freshSampleRawRecommendation());
    savingRecommendation = new RecommendationExtra(
      freshSavingRawRecommendation()
    );
    performanceRecommendation = new RecommendationExtra(
      freshPerformanceRawRecommendation()
    );
    fakeStore = rootStoreFactory();
    fakeStore.commit("recommendationsStore/addRecommendation", recommendation);
    fakeStore.commit(
      "recommendationsStore/addRecommendation",
      savingRecommendation
    );
    fakeStore.commit(
      "recommendationsStore/addRecommendation",
      performanceRecommendation
    );

    // So the warning about not finding data-app doesn't happen
    const app = document.createElement("div");
    app.setAttribute("data-app", "true");
    document.body.appendChild(app);
  });

  it("When there are no selected recommendations, there is no footer", async () => {
    const wrapper = mount(Footer, { store: fakeStore, vuetify: vuetify });
    const selected = new Array<RecommendationExtra>();
    selected.push(recommendation);
    fakeStore.commit("coreTableStore/setSelected", selected);
    await wrapper.vm.$nextTick();

    expect(wrapper.findComponent({ ref: "footer" }).exists()).toBe(true);
  });

  it("When there are selected recommendations, there is a footer", async () => {
    const wrapper = mount(Footer, { store: fakeStore, vuetify: vuetify });
    const selected = new Array<RecommendationExtra>();
    selected.push(recommendation);
    fakeStore.commit("coreTableStore/setSelected", selected);
    await wrapper.vm.$nextTick();

    expect(wrapper.findComponent({ ref: "footer" }).exists()).toBe(true);
  });

  it("When there are selected recommendations and they become unselected, there is no footer", async () => {
    const wrapper = mount(Footer, { store: fakeStore, vuetify: vuetify });
    const selected = new Array<RecommendationExtra>();
    selected.push(recommendation);
    fakeStore.commit("coreTableStore/setSelected", selected);
    await wrapper.vm.$nextTick(); // Needed since Vue 2.0 to have reactivity

    fakeStore.commit(
      "coreTableStore/setSelected",
      new Array<RecommendationExtra>()
    );
    await wrapper.vm.$nextTick();

    expect(wrapper.findComponent({ ref: "footer" }).exists()).toBe(false);
  });

  it("When saving recommendation is selected the text on the footer is correct", async () => {
    const wrapper = mount(Footer, { store: fakeStore, vuetify: vuetify });
    const selected = new Array<RecommendationExtra>();
    selected.push(savingRecommendation);
    fakeStore.commit("coreTableStore/setSelected", selected);
    await wrapper.vm.$nextTick(); // Needed since Vue 2.0 to have reactivity

    expect(
      wrapper
        .findComponent({ ref: "footer" })
        .find("[data-name=footer-summary]")
        .text()
    ).toEqual(
      `Apply 1 recommendation(s):  Use less resources and save ${-savingRecommendation.costCol.toFixed(
        2
      )}$ each week.`
    );
  });

  it("When performance recommendation is selected the text on the footer is correct", async () => {
    const wrapper = mount(Footer, { store: fakeStore, vuetify: vuetify });
    const selected = new Array<RecommendationExtra>();
    selected.push(performanceRecommendation);
    fakeStore.commit("coreTableStore/setSelected", selected);
    await wrapper.vm.$nextTick(); // Needed since Vue 2.0 to have reactivity

    expect(
      wrapper
        .findComponent({ ref: "footer" })
        .find("[data-name=footer-summary]")
        .text()
    ).toEqual(
      `Apply 1 recommendation(s):  Increase performance of 1 machine by spending ${performanceRecommendation.costCol.toFixed(
        2
      )}$ more each week.`
    );
  });

  it("When both kinds of recommendations are selected the text on the footer is correct", async () => {
    const wrapper = mount(Footer, { store: fakeStore, vuetify: vuetify });
    const selected = new Array<RecommendationExtra>();
    selected.push(performanceRecommendation);
    selected.push(savingRecommendation);
    fakeStore.commit("coreTableStore/setSelected", selected);
    await wrapper.vm.$nextTick(); // Needed since Vue 2.0 to have reactivity

    expect(
      wrapper
        .findComponent({ ref: "footer" })
        .find("[data-name=footer-summary]")
        .text()
    ).toEqual(
      `Apply 2 recommendation(s):  Increase performance of 1 machine by spending ${performanceRecommendation.costCol.toFixed(
        2
      )}$ more each week. Use less resources and save ${-savingRecommendation.costCol.toFixed(
        2
      )}$ each week.`
    );
  });

  it("There is a button in the footer", async () => {
    const wrapper = mount(Footer, { store: fakeStore, vuetify: vuetify });
    const selected = new Array<RecommendationExtra>();
    selected.push(performanceRecommendation);
    selected.push(savingRecommendation);
    fakeStore.commit("coreTableStore/setSelected", selected);
    await wrapper.vm.$nextTick(); // Needed since Vue 2.0 to have reactivity

    expect(
      wrapper
        .findComponent({ ref: "footer" })
        .find("[data-name=footer-button]")
        .exists()
    ).toBe(true);
  });

  it("Dialog window is not visible before clicking the button", async () => {
    const wrapper = mount(Footer, { store: fakeStore, vuetify: vuetify });
    const selected = new Array<RecommendationExtra>();
    selected.push(performanceRecommendation);
    selected.push(savingRecommendation);
    fakeStore.commit("coreTableStore/setSelected", selected);
    await wrapper.vm.$nextTick(); // Needed since Vue 2.0 to have reactivity

    expect(
      wrapper
        .findComponent({ ref: "dialog" })
        .find("[data-name=dialog]")
        .exists()
    ).toBe(false);
  });

  it("After clicking the button a dialog window opens", async () => {
    const wrapper = mount(Footer, { store: fakeStore, vuetify: vuetify });
    const selected = new Array<RecommendationExtra>();
    selected.push(performanceRecommendation);
    selected.push(savingRecommendation);
    fakeStore.commit("coreTableStore/setSelected", selected);
    await wrapper.vm.$nextTick(); // Needed since Vue 2.0 to have reactivity

    await wrapper
      .findComponent({ ref: "footer" })
      .find("[data-name=footer-button]")
      .trigger("click");

    await wrapper.vm.$nextTick();

    expect(
      wrapper.findComponent({ ref: "dialog" }).find("[data-name=dialog]")
        .element
    ).toBeVisible();
  });

  it("After canceling applying the dialog window closes and recommendations remain selected", async () => {
    const wrapper = mount(Footer, { store: fakeStore, vuetify: vuetify });
    const selected = new Array<RecommendationExtra>();
    selected.push(performanceRecommendation);
    selected.push(savingRecommendation);
    fakeStore.commit("coreTableStore/setSelected", selected);
    await wrapper.vm.$nextTick(); // Needed since Vue 2.0 to have reactivity

    await wrapper
      .findComponent({ ref: "footer" })
      .find("[data-name=footer-button]")
      .trigger("click");

    await wrapper.vm.$nextTick();

    await wrapper
      .findComponent({ ref: "dialog" })
      .find("[data-name=cancel-button]")
      .trigger("click");

    await wrapper.vm.$nextTick();

    expect(
      wrapper.findComponent({ ref: "dialog" }).find("[data-name=dialog]")
        .element
    ).not.toBeVisible();

    expect(performanceRecommendation.statusCol === "ACTIVE");

    expect(fakeStore.state.coreTableStore?.selected).toEqual(selected);
  });

  it("After clicking yes, the dialog window closes and recommendations change status from active", async () => {
    const wrapper = mount(Footer, { store: fakeStore, vuetify: vuetify });
    const selected = new Array<RecommendationExtra>();
    selected.push(performanceRecommendation);
    selected.push(savingRecommendation);
    fakeStore.commit("coreTableStore/setSelected", selected);
    await wrapper.vm.$nextTick(); // Needed since Vue 2.0 to have reactivity

    await wrapper
      .findComponent({ ref: "footer" })
      .find("[data-name=footer-button]")
      .trigger("click");

    await wrapper.vm.$nextTick();

    await wrapper
      .findComponent({ ref: "dialog" })
      .find("[data-name=yes-button]")
      .trigger("click");

    await wrapper.vm.$nextTick();

    expect(
      wrapper.findComponent({ ref: "dialog" }).find("[data-name=dialog]")
        .element
    ).not.toBeVisible();

    expect(fakeStore.state.coreTableStore?.selected.length).toEqual(0);
    expect(performanceRecommendation.statusCol !== "ACTIVE");
    expect(savingRecommendation.statusCol !== "ACTIVE");
  });
});
