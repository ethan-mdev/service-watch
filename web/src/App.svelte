<script>
  import TopBar from './lib/components/TopBar.svelte';
  import ResourceSection from './lib/components/ResourceSection.svelte';
  import GraphsSection from './lib/components/GraphsSection.svelte';
  import BottomSection from './lib/components/BottomSection.svelte';
  import Footer from './lib/components/Footer.svelte';
</script>

<!-- Main container -->
<div class="bg-black min-h-screen flex flex-col items-center">
  <!-- Content area -->
  <div class="max-w-[80%] bg-white flex grow min-w-[80%] flex-col">
    <TopBar />
    <ResourceSection />
    <GraphsSection />
    <BottomSection />
    <Footer />
  </div>
</div>