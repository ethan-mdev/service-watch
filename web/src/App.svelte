<script>
  import TopBar from './lib/components/TopBar.svelte';
  import ResourceSection from './lib/components/ResourceSection.svelte';
  import GraphsSection from './lib/components/GraphsSection.svelte';
  import BottomSection from './lib/components/BottomSection.svelte';
  import Footer from './lib/components/Footer.svelte';
</script>

<div class="bg-black text-neutral-200 antialiased min-h-screen">
  <TopBar />
  <main class="mx-auto max-w-7xl px-6 sm:px-8 md:px-12 lg:px-16 xl:px-4 py-6 space-y-8">
    <ResourceSection />
    <GraphsSection />
    <BottomSection />
  </main>

  <Footer />
</div>