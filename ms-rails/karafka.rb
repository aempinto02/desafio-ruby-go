# frozen_string_literal: true

puts "Loading Karafka configuration..."

require 'karafka'
require_relative 'app/consumers/go_topic_consumer'

class KarafkaApp < Karafka::App
  setup do |config|
    config.consumer_persistence = !Rails.env.development?
    config.kafka = { 'bootstrap.servers': 'kafka:29092' }
    config.client_id = 'ms-rails'
    
    consumer_groups.draw do
      topic :go_to_rails do
        consumer GoTopicConsumer
      end
      puts "Subscribed to topic: go_to_rails"
    end

    config.admin do |admin|
      admin.kafka['allow.auto.create.topics'] = true
    end

    Karafka.logger.info("Karafka configuration: #{config.inspect}")
  end

  Karafka.monitor.subscribe(Karafka::Instrumentation::LoggerListener.new)

  Karafka.producer.monitor.subscribe(
    WaterDrop::Instrumentation::LoggerListener.new(
      Karafka.logger,
      log_messages: false
    )
  )
end

KarafkaApp.run!
puts "KarafkaApp booted successfully"
